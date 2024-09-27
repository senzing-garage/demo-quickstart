#!/usr/bin/env python
# coding: utf-8

# # Senzing Hello World
#
# The following is a quick test of connectivity to Senzing using Python over gRPC.
#

# In[ ]:


import json

import grpc
from senzing_grpc import SzAbstractFactory

sz_abstract_factory = SzAbstractFactory(
    grpc_channel=grpc.insecure_channel("localhost:8261")
)
sz_product = sz_abstract_factory.create_sz_product()
print(json.dumps(json.loads(sz_product.get_version()), indent=2))
