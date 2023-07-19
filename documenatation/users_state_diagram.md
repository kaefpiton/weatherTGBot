stateDiagram-v2
    [*] --> authorised
    authorised --> admin_auth
    authorised --> get_weather
    admin_auth --> admin_enter
    admin_auth --> authorised
    admin_enter --> authorised