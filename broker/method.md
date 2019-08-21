
MY IMPLEMENTATION
connection - class 10
    method:
        ConnectionStart    - 10
        ConnectionStartOk  - 11
        ConnectionOpen     - 20
        ConnectionOpenOk   - 21
        ConnectionClose    - 30
        ConnectionCloseOk  - 31
channel - class 20
    method:
        ChannelOpen - 10
        ChannelOpenOk - 11
        ChannelFlow - 20
        ChannelFlowOk - 21
        ChannelClose - 30
        ChannelCloseOk - 31
exchange - class 30
    method:
        ExchangeDeclare - 10
        ExchangeDeclareOk  - 11
        ExchangeDelete - 20
        ExchangeDeleteOk - 21
        ExchangeBind  - 30
        ExchangeBindOk - 31
        ExchangeUnbind - 40
        ExchangeUnbindOk - 41
queue - class 40
    method:
        QueueDeclare - 10
        QueueDeclareOk - 11
        QueueBind      - 20
        QueueBindOk    - 21
        QueueUnbind    - 30
        QueueUnbindOk  - 31
        QueueDelete    - 40
        QueueDeleteOk  - 41
Qos - class 50
    method:
        BasicConsume - 10
        BasicConsumeOk - 11
        BasicCancel - 20
        BasicCancelOk - 21
        BasicPublish - 30
        BasicReturn  - 40
        BasicDeliver - 50
        BasicAck     - 60
        BasicNack    - 70

***************************************

connection - 10 class
    method:
        connectionStart    - 10
        connectionStartOk  - 11
        connectionSecure   - 20
        connectionSecureOk - 21
        connectionTune     - 30
        connectionTuneOk   - 31
        connectionOpen     - 40
        connectionOpenOk   - 41
        connectionClose    - 50
        connectionCloseOk  - 51
        connectionBlocked  - 60
        connectionUnblocked- 61
channel - 20 class
    method:
        channelOpen - 10
        channelOpenOk - 11
        channelFlow - 20
        channelFlowOk - 21
        channelClose - 30
        channelCloseOk - 31
exchange - 40 class
    method:
        exchangeDeclare - 10
        exchangeDeclareOk  - 11
        exchangeDelete - 20
        exchangeDeleteOk - 21
        exchangeBind  - 30
        exchangeBindOk - 31
        exchangeUnbind - 40
        exchangeUnbindOk - 41
queue - 50 class
    method:
        queueDeclare - 10
        queueDeclareOk - 11
        queueBind      - 20
        queueBindOk    - 21
        queueUnbind    - 50
        queueUnbindOk  - 51
        queuePurge     - 30
        queuePurgeOk   - 31
        queueDelete    - 40
        queueDeleteOk  - 41
Qos - 60 class
    method:
        basicQos - 10
        basicQosOk - 11
        basicConsume - 20
        basicConsumeOk - 21
        basicCancel - 30
        basicCancelOk - 31
        basicPublish - 40
        basicReturn  - 50
        basicDeliver - 60
        basicGet     - 70
        basicGetOk   - 71
        basicGetEmpty - 72
        basicAck      - 80
        basicReject   - 90
        basicRecoverAsync - 100


