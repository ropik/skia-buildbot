AlertServer
===========

AlertServer is a server which periodically queries InfluxDB and generates
alerts based on rules defined in alerts.cfg.


### AlertServer ###
It needs the following metadata values:

    cookiesalt
    client_id
    client_secret
    influxdb_name
    influxdb_password
    gmail_clientid
    gmail_clientsecret
    gmail_cached_token

The client_id and client_secret come from here:

    https://console.developers.google.com/project/31977622648/apiui/credential

Look for the Client ID that has a Redirect URI for skiamonitor.com.

For 'cookiesalt' and the influx db values search for 'skiamonitor' in valentine.

The gmail_clientid and gmail_clientsecret come from here:

    https://console.developers.google.com/project/31977622648/apiui/credential

Look for the section titled, "Client ID for native application."

The gmail_cached_token can be generated by running the server and clicking the
authorization link while signed in as skia.buildbots@gmail.com