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
        channelClose - 40
        channelCloseOk - 41
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


