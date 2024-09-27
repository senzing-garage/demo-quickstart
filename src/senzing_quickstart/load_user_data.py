#!/usr/bin/env python
# coding: utf-8

# # Load user data
#
# These instructions show how to load user data into the Senzing engine.

# ## Upload file of records

# The following instructions show how to upload a file of JSON lines.
# Here is an [example file].
#
# 1. In upper-left corner of JuypterLab, click on the
#    ![Upload Files](img/upload-icon.png)
#    ("Upload Files") icon.
# 1. Choose the file of JSON lines to upload.
#
# [example file]: https://raw.githubusercontent.com/senzing-garage/demo-quickstart/refs/heads/main/testdata/example-data-for-senzing.json

# ## Prepare Python enviroment

# Import Python packages.

# In[ ]:


import json

import grpc
from senzing_grpc import SzAbstractFactory, SzEngineFlags, SzError

# Set environment specific variables.

# In[ ]:


home_path = "/notebooks/"


# ![Modify](img/pencil2.png) **Modify the following.** Identify uploaded file.

# In[ ]:


filepath = home_path + "example-data-for-senzing.json"


# ## Identify data sources

# Discover `DATA_SOURCE` values in records.

# In[ ]:


datasources = []

with open(filepath, "r", encoding="utf-8") as file:
    for line in file:
        line_as_dict = json.loads(line)
        datasource = line_as_dict.get("DATA_SOURCE")
        if datasource not in datasources:
            datasources.append(datasource)

print(f"Found the following DATA_SOURCE values in the data: {datasources}")


# ## Update Senzing configuration

# Create an abstract factory for accessing Senzing via gRPC.

# In[ ]:


sz_abstract_factory = SzAbstractFactory(
    grpc_channel=grpc.insecure_channel("localhost:8261")
)


# Create Senzing objects.

# In[ ]:


sz_config = sz_abstract_factory.create_sz_config()
sz_configmanager = sz_abstract_factory.create_sz_configmanager()
sz_diagnostic = sz_abstract_factory.create_sz_diagnostic()
sz_engine = sz_abstract_factory.create_sz_engine()


# Get current Senzing configuration.

# In[ ]:


old_config_id = sz_configmanager.get_default_config_id()
old_json_config = sz_configmanager.get_config(old_config_id)
config_handle = sz_config.import_config(old_json_config)


# Add DataSources to Senzing configuration.

# In[ ]:


for datasource in datasources:
    try:
        sz_config.add_data_source(config_handle, datasource)
    except SzError as err:
        print(err)


# Persist new Senzing configuration.

# In[ ]:


new_json_config = sz_config.export_config(config_handle)
new_config_id = sz_configmanager.add_config(new_json_config, "Add user datasources")
sz_configmanager.replace_default_config_id(old_config_id, new_config_id)


# With the change in Senzing configuration, Senzing objects need to be updated.

# In[ ]:


sz_engine.reinitialize(new_config_id)
sz_diagnostic.reinitialize(new_config_id)


# ## Add records

# Call Senzing to add records.

# In[ ]:


with open(filepath, "r") as file:
    for line in file:
        try:
            line_as_dict = json.loads(line)
            info = sz_engine.add_record(
                line_as_dict.get("DATA_SOURCE"),
                line_as_dict.get("RECORD_ID"),
                line,
                SzEngineFlags.SZ_WITH_INFO,
            )
            print(info)
        except SzError as err:
            print(err)
