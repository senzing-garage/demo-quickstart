#!/usr/bin/env python
# coding: utf-8

# # Show method help

# Import Python packages.

# In[ ]:


import grpc

from senzing_grpc import SzAbstractFactory


# Create an abstract factory for accessing Senzing via gRPC.

# In[ ]:


sz_abstract_factory = SzAbstractFactory(
    grpc_channel=grpc.insecure_channel("localhost:8261")
)


# Create Senzing object.

# In[ ]:


sz_engine = sz_abstract_factory.create_sz_engine()


# List all methods for a Senzing object.

# In[ ]:


print(sz_engine.help())


# Print help for a specific method.

# In[ ]:


print(sz_engine.help("get_entity_by_record_id"))

