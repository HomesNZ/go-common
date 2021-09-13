# Settings

Example of serverless yml file:
```
resources:
  Resources:
    MainQueue:
      Type: AWS::SQS::Queue
      Properties:
        QueueName: ${self:provider.environment.SQS_QUOTES}
        ReceiveMessageWaitTimeSeconds: 10
        VisibilityTimeout: 1800 # 30 mins -- prevent other consumers from processing the message again
        RedrivePolicy:
          deadLetterTargetArn:
            Fn::GetAtt:
              - "DeadLetterQueue"
              - "Arn"
          maxReceiveCount: 5 # as per AWS Documentation.
    DeadLetterQueue:
      Type: AWS::SQS::Queue
      Properties:
        QueueName: ${self:provider.environment.SQS_QUOTES}_dead_letter
        MessageRetentionPeriod: 1209600 # 14 days in seconds
```