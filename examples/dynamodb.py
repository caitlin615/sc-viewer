# /webhook consumes from an AWS Lambda
# the Lambda function listens for item creations in a DynamoDB database
# and send the data to the webhook. To use this, create a python Lambda function
# and add this script. Add a trigger to listen to a DynamoDB event stream
# WEBHOOK_URL and WEBHOOK_SECRET should be set in the Lambda function

from __future__ import print_function

import json
import urllib
import boto3
import os
import logging
from datetime import datetime
import base64

from urllib2 import Request, urlopen, URLError, HTTPError

logging.info("Start logger")
log = logging.getLogger()
log.setLevel(logging.INFO)

webhook_url = os.getenv("WEBHOOK_URL")
secret = os.getenv("WEBHOOK_SECRET")

def lambda_handler(event, context):
    #print("Received event: " + json.dumps(event, indent=2))
    items = []
    for record in event['Records']:
        print(record['eventID'])
        if record['eventName'] != 'MODIFY':
            continue
       item = {}
       date_sent = record['dynamodb']['NewImage'].get('date_sent')
       if date_sent:
           item["date_sent"] = date_sent
        item["id"] = record['dynamodb']['NewImage']['id']['S']
        item["shareMode"] = record['dynamodb']['NewImage']['exporterShareMode']['S']
        urlRecord = record['dynamodb']['NewImage']['urls']['M']
        item["urls"] = {}
        for mimeType, value in urlRecord.iteritems():
            item["urls"][mimeType] = value['S']

        items.append(item)
    try:
        req = Request(webhook_url, json.dumps(items))
        base64string = base64.b64encode('%s:%s' % ("lambda", secret))
        req.add_header("Authorization", "Basic %s" % base64string)
        try:
            response = urlopen(req)
            response.read()
            # log.info("Success:", items)
        except HTTPError as e:
            log.error("Request failed: %d %s", e.code, e.reason)
            raise e
        except URLError as e:
            log.error("Server connection failed: %s", e.reason)
            raise e
    except Exception as e:
        log.info(e)
        raise e
