#!/usr/bin/env python
# coding: utf-8

# # Load Senzing truth-sets
# 
# These instructions load the [Senzing truth-sets] into the Senzing engine.
# 
# [Senzing truth-sets]: https://github.com/Senzing/truth-sets

# ## Prepare Python enviroment

# Import Python packages.

# In[ ]:


import json
import shutil

import grpc
import requests
from senzing_grpc import SzAbstractFactory, SzEngineFlags, SzError


# Set environment specific variables.

# In[ ]:


home_path = "/notebooks/"
truth_set_url_prefix = "https://raw.githubusercontent.com/Senzing/truth-sets/refs/heads/main/truthsets/demo/"
truth_set_filenames = [
    "customers.json",
    "reference.json",
    "watchlist.json",
]


# Download truth-set files.

# In[ ]:


for filename in truth_set_filenames:
    url = truth_set_url_prefix + filename
    filepath = home_path + filename
    response = requests.get(url, stream=True, timeout=10)
    response.raw.decode_content = True
    with open(filepath, "wb") as file:
        shutil.copyfileobj(response.raw, file)


# ## Identify data sources

# Discover `DATA_SOURCE` values in records.

# In[ ]:


datasources = []

for filename in truth_set_filenames:
    filepath = home_path + filename
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


sz_abstract_factory = SzAbstractFactory(
    grpc_channel=grpc.insecure_channel("localhost:8261")
)
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
new_config_id = sz_configmanager.add_config(
    new_json_config, "Add TruthSet datasources"
)
sz_configmanager.replace_default_config_id(old_config_id, new_config_id)    


# With the change in Senzing configuration, Senzing objects need to be updated.

# In[ ]:


sz_engine.reinitialize(new_config_id)
sz_diagnostic.reinitialize(new_config_id)


# ## Add records

# Call Senzing to add records.

# In[ ]:


for filename in truth_set_filenames:
    filepath = home_path + filename
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


# ## View results

# Retrieve an entity by identifying a record of the entity.

# In[ ]:


customer_1070_entity = sz_engine.get_entity_by_record_id("CUSTOMERS", "1070")
print(json.dumps(json.loads(customer_1070_entity), indent=2))


# Search for entities by attributes.

# In[ ]:


search_query = {
    "name_full": "robert smith",
    "date_of_birth": "11/12/1978",
}
search_result = sz_engine.search_by_attributes(json.dumps(search_query))
print(json.dumps(json.loads(search_result), indent=2))

