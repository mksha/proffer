# copyami resource properties

> **name:**
> - **type:** String
> - **required:** True
> - **allowed-values:** Any valid string.

<hr style="border:2px solid gray"> </hr>

> **type:**
> - **type:** String
> - **required:** True
> - **allowed-values:** aws-copyami

<hr style="border:2px solid gray"> </hr>

> **config:**
> - **type:** Dict
> - **required:** True
> >
> **Properties/Keys:**
>> **source:**
>> - **type:** Dict
>> - **required:** True
>> >
>> **Properties/Keys:**
>>> **profile:**
>>> - **type:** String 
>>> - **required:** Optional
>>> - **allowed-values:** Valid aws profile name
>>> - **description:** Aws Profile to get creds for source account.
> > 
>> <hr style="border:2px solid gray"> </hr>
> >
>>> **roleArn:**
>>> - **type:** String
>>> - **required:** Optional
>>> - **allowed-values:** valid aws role arn
>>> - **description:** AWS Role ARN to get creds for source account.
> >
>> <hr style="border:2px solid gray"> </hr>
> >
>>> **region:**
>>> - **type:** String
>>> - **required:** True
>>> - **allowed-values:** Valid aws region.
>>> - **description:** Source ami region.
> >
>> <hr style="border:2px solid gray"> </hr>
> >
>>> **amiFilters:**
>>> - **type:** Dict.
>>> - **required:** True
>>> - **allowed-values:** Valid AWS ami filters.
>>> - **description:** AMI filters to uniquely identify source ami.
>
> <hr style="border:2px solid gray"> </hr>
>
>> **target:**
>> - **type:** Dict
>> - **required:** True
>> >
>> **Properties/Keys:**
>>> **regions:**
>>> - **type:** List
>>> - **required:** True
>>> - **allowed-values:** List of valid aws regions.
>>> - **description:** Target AWS regions to which we want to copy the source ami.
> >
>>> **copyTagsAcrossRegions:**
>>> - **type:** Boolean
>>> - **required:** Optional
>>> - **allowed-values:** [true, false].
>>> - **description:** Set this flag to true if you want to copy the source ami tags to target ami.
> >
>>> **addExtraTags:**
>>> - **type:** Dict
>>> - **required:** Optional
>>> - **allowed-values:** Valid AWS tags.
>>> - **description:** Add extra tags to target ami in the form of `key:value` .
