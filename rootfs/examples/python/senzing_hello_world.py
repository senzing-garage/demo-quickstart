#!/usr/bin/env python3

import json

import grpc
from senzing_grpc import SzAbstractFactory, SzAbstractFactoryParameters

FACTORY_PARAMETERS: SzAbstractFactoryParameters = {
    "grpc_channel": grpc.insecure_channel("localhost:8261"),
}

sz_abstract_factory = SzAbstractFactory(**FACTORY_PARAMETERS)
sz_product = sz_abstract_factory.create_product()
print(json.dumps(json.loads(sz_product.get_version()), indent=2))
