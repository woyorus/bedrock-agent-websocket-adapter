# Amazon Bedrock Agent adapter for WebSockets

This is used as example: https://medium.com/@rachmanarief/real-time-chatbots-with-amazon-bedrock-agent-websocket-and-golang-6ee73e4cd838

Create `.env` file with following content:

```
AGENT_ID=YOUR_AGENT_ID
AGENT_ALIAS_ID=YOUR_ALIAS_ID
AWS_SHARED_CONFIG_PROFLE=default
```

(replace YOUR_AGENT_ID and YOUR_ALIAS_ID with your values from Agent setup).

AWS Creds are configured as usual with aws-sdk-go - see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

Then connect as WebSocket target to `localhost:7001/ws`  
Send messages and receive answers from Bedrock Agent.
