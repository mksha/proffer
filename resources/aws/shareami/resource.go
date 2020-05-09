package shareami

import (
	"fmt"
	"log"
	"os"
	"strconv"

	clog "example.com/proffer/common/clogger"
	awscommon "example.com/proffer/resources/aws/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

var (
	clogger = clog.New(os.Stdout, "aws-shareami | ", log.Lmsgprefix)
)

type RawSrcAmiInfo struct {
	Profile    *string             `yaml:"profile"`
	RoleArn    *string             `yaml:"roleArn"`
	AmiFilters map[*string]*string `yaml:"amiFilters"`
}

type SrcAmiInfo struct {
	CredsInfo       map[string]string
	AccountID       *string
	Filters         []*ec2.Filter
	RegionAmiErrMap RegionAmiErrMap
}

type RawAccountRegionMapping struct {
	AccountID                 int                 `yaml:"accountId"`
	Profile                   *string             `yaml:"profile"`
	RoleArn                   *string             `yaml:"roleArn"`
	Regions                   []*string           `yaml:"regions"`
	AddExtraTags              map[*string]*string `yaml:"addExtraTags"`
	CopyTagsAcrossAccounts    bool                `yaml:"copyTagsAcrossAccounts"`
	AddCreateVolumePermission bool                `yaml:"addCreateVolumePermission"`
}

type AccountRegionMapping struct {
	CopyTags  bool
	AddCVP    bool
	Tags      []*ec2.Tag
	CredsInfo map[string]string
	AccountID *string
	Image     *ec2.Image
	Regions   []*string
}

type Target struct {
	AccountRegionMappingList    []RawAccountRegionMapping `yaml:"accountRegionMappingList"`
	CopyTagsAcrossAccounts      bool                      `yaml:"copyTagsAcrossAccounts"`
	CommonRegions               []*string                 `yaml:"commonRegions"`
	AddCreateVolumePermission   bool                      `yaml:"addCreateVolumePermission"`
	ModAccountRegionMappingList []AccountRegionMapping    `yaml:"-"`
}

type Config struct {
	Source     RawSrcAmiInfo `yaml:"source"`
	Target     Target        `yaml:"target"`
	SrcAmiInfo SrcAmiInfo    `yaml:"-"`
}

type Resource struct {
	Config Config `yaml:"config"`
}

type (
	AmiInfo struct {
		Ami   *ec2.Image
		Error error
	}
	RegionAmiErrMap map[*string]AmiInfo
	Flag struct{
		AddCVP bool
	}
	AccountFlagMap  map[*string]Flag
)

func (r *Resource) Prepare(rawConfig map[string]interface{}) error {
	var c Config

	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	clogger.Warn("Gathering Information...")

	if err := mapstructure.Decode(rawConfig, &c); err != nil {
		return err
	}

	r.Config = c

	r.Config.SrcAmiInfo = prepareSrcAmiInfo(r.Config.Source)

	sess, err := awscommon.GetAwsSession(r.Config.SrcAmiInfo.CredsInfo)
	if err != nil {
		return err
	}

	accountInfo, err := awscommon.GetAccountInfo(sess)
	if err != nil {
		return err
	}

	r.Config.SrcAmiInfo.AccountID = accountInfo.Account

	regions := r.Config.Target.getTargetRegions()

	if err := r.Config.SrcAmiInfo.prepareTargetRegionAmiMapping(regions); err != nil {
		for region, amiInfo := range r.Config.SrcAmiInfo.RegionAmiErrMap {
			if amiInfo.Error != nil {
				clogger.Infof("Source AMI Not Found In Account: %s Region: %s", *r.Config.SrcAmiInfo.AccountID, *region)
				clogger.Error(amiInfo.Error)
			}
		}

		return fmt.Errorf("failed: To Get Source AMI Information For Some Region(s)")
	}

	r.prepareAccountRegionMappingList()
	r.Config.Target.setCommonPropertiesIfAny()

	clogger.Success("Successfully Gathered All Info Needed For Source")
	clogger.Info("")

	return nil
}

func (r *Resource) Run() error {

	if err := r.Config.apply(); err != nil {
		return err
	}

	return nil
}

func prepareSrcAmiInfo(rawSrcAmiInfo RawSrcAmiInfo) SrcAmiInfo {
	var amiFilters []*ec2.Filter

	for filterName, filterValue := range rawSrcAmiInfo.AmiFilters {
		f := &ec2.Filter{
			Name:   filterName,
			Values: []*string{filterValue},
		}
		amiFilters = append(amiFilters, f)
	}

	srcAmiInfo := SrcAmiInfo{
		Filters:         amiFilters,
		CredsInfo:       make(map[string]string, 2),
		RegionAmiErrMap: make(RegionAmiErrMap, 0),
	}

	if rawSrcAmiInfo.RoleArn != nil {
		srcAmiInfo.CredsInfo["getCredsUsing"] = "roleArn"
		srcAmiInfo.CredsInfo["roleArn"] = *rawSrcAmiInfo.RoleArn
	} else if rawSrcAmiInfo.Profile != nil {
		srcAmiInfo.CredsInfo["getCredsUsing"] = "profile"
		srcAmiInfo.CredsInfo["profile"] = *rawSrcAmiInfo.Profile
	}

	return srcAmiInfo
}

func (r *Resource) prepareAccountRegionMappingList() {
	accountRegionMappingList := make([]AccountRegionMapping, 0)

	for _, rawAccountRegionMapping := range r.Config.Target.AccountRegionMappingList {
		accountRegionMapping := AccountRegionMapping{
			CopyTags:  rawAccountRegionMapping.CopyTagsAcrossAccounts,
			Tags:      awscommon.FormEc2Tags(rawAccountRegionMapping.AddExtraTags),
			Regions:   rawAccountRegionMapping.Regions,
			AccountID: aws.String(strconv.Itoa(rawAccountRegionMapping.AccountID)),
			CredsInfo: make(map[string]string, 2),
		}

		if rawAccountRegionMapping.RoleArn != nil {
			accountRegionMapping.CredsInfo["getCredsUsing"] = "roleArn"
			accountRegionMapping.CredsInfo["roleArn"] = *rawAccountRegionMapping.RoleArn
		} else if rawAccountRegionMapping.Profile != nil {
			accountRegionMapping.CredsInfo["getCredsUsing"] = "profile"
			accountRegionMapping.CredsInfo["profile"] = *rawAccountRegionMapping.Profile
		}

		accountRegionMappingList = append(accountRegionMappingList, accountRegionMapping)
	}

	r.Config.Target.ModAccountRegionMappingList = accountRegionMappingList
}

func (t *Target) setCommonPropertiesIfAny() {
	if t.CopyTagsAcrossAccounts {
		for i := 0; i < len(t.ModAccountRegionMappingList); i++ {
			t.ModAccountRegionMappingList[i].CopyTags = t.CopyTagsAcrossAccounts
		}
	}

	if t.AddCreateVolumePermission {
		for i := 0; i < len(t.ModAccountRegionMappingList); i++ {
			t.ModAccountRegionMappingList[i].AddCVP = t.AddCreateVolumePermission
		}
	}
}

func (sai *SrcAmiInfo) prepareTargetRegionAmiMapping(regions []*string) (err error) {
	regionAmiErrChan := make(chan RegionAmiErrMap)
	defer close(regionAmiErrChan)

	for _, targetRegion := range regions {
		sai := *sai
		go prepareRegionAmiErrMap(sai, targetRegion, regionAmiErrChan)
	}

	for i := 0; i < len(regions); i++ {
		regionAmiInfo := <-regionAmiErrChan
		for region, amiInfo := range regionAmiInfo {
			sai.RegionAmiErrMap[region] = amiInfo

			if amiInfo.Error != nil {
				err = amiInfo.Error
			}
		}
	}

	return
}

func prepareRegionAmiErrMap(sai SrcAmiInfo, region *string, regionAmiErrChan chan<- RegionAmiErrMap) {
	regionAmiErr := make(RegionAmiErrMap, 0)

	sess, err := awscommon.GetAwsSession(sai.CredsInfo)
	if err != nil {
		regionAmiErr[region] = AmiInfo{Error: err}
		regionAmiErrChan <- regionAmiErr

		return
	}

	sess.Config.Region = region
	images, err := awscommon.GetAmiInfo(sess, sai.Filters)

	if err != nil {
		regionAmiErr[region] = AmiInfo{Error: err}

		regionAmiErrChan <- regionAmiErr

		return
	}

	image := images[0]

	clogger.Infof("Source AMI: %s Found In Account: %s In Region: %s", *image.Name, *sai.AccountID, *region)

	regionAmiErr[region] = AmiInfo{Ami: image}

	regionAmiErrChan <- regionAmiErr
}

func (t Target) getTargetRegions() []*string {
	regions := make([]*string, 0)

	for _, rawAccountRegionMapping := range t.AccountRegionMappingList {
		regions = append(regions, rawAccountRegionMapping.Regions...)
	}

	regions = append(regions, t.CommonRegions...)

	return regions
}

func (t Target) getTargetAccountFlagMap() AccountFlagMap {
	accountFlagMap := make(AccountFlagMap, 0)

	for _, accountRegionMapping := range t.ModAccountRegionMappingList {
		account := accountRegionMapping.AccountID
		accountFlagMap[account] = Flag{AddCVP: accountRegionMapping.AddCVP}
	}

	return accountFlagMap
}
