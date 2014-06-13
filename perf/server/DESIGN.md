DESIGN
======


Overview
--------
Provides interactive dashboard for Skia performance data.


Architecture
------------

Perf Stats Database
-------------------

Annotations Database
--------------------

A Cloud SQL (a cloud version of MySQL) database is used to keep information on
Skia git revisions and their corresponding annotations. The database will be
updated when users add/edit/delete annotations via the dashboard UI.

All passwords for MySQL are stored in valentine (search "skia perf").

To connect to the database from authorized network (including skia-perf GCE):

    $ mysql -h 173.194.104.24 -u root -p

Initial setup of the database, the users, and the tables:

    CREATE DATABASE skia;
    USE skia;
    CREATE USER 'readonly'@'%' IDENTIFIED BY <password in valentine>;
    GRANT SELECT ON *.* TO 'readonly'@'%';
    CREATE USER 'readwrite'@'%' IDENTIFIED BY <password in valentine>;
    GRANT SELECT, DELETE, UPDATE, INSERT ON *.* TO 'readwrite'@'%';

    // Table for storing annotations.
    CREATE TABLE notes (
      id     INT           NOT NULL AUTO_INCREMENT PRIMARY KEY,
      type   TINYINT,
      author VARCHAR(40),
      notes  VARCHAR(200)  NOT NULL,
      UNIQUE INDEX (id)
    );

    // Table for storing git revision information.
    CREATE TABLE githash (
      ts        INT           NOT NULL PRIMARY KEY,
      gitnumber INT           NOT NULL,
      githash   VARCHAR(40)   NOT NULL,
      author    VARCHAR(40)   NOT NULL,
      message   VARCHAR(200)  NOT NULL,
      UNIQUE INDEX (ts)
    );

    // Table for mapping revisions and annotations. This support many-to-many
    // mapping.
    CREATE TABLE githashnotes (
      ts INT   NOT NULL,
      id INT   NOT NULL,

      FOREIGN KEY (ts) REFERENCES githash(ts),
      FOREIGN KEY (id) REFERENCES notes(id),
      UNIQUE INDEX (ts, id)
    );

Common queries that the dashboard will use:

    INSERT INTO notes (type, author, notes) VALUES (1, 'bsalomon', 'Alert!');

    SELECT LAST_INSERT_ID();

    INSERT INTO githashnotes (ts, id) VALUES (<githash_ts>, <last_insert_id>);

The above set of commands will usually be used together to add new annotations
and associate them with corresponding git commits. The commands below remove an
annotation and its associations with any commit.

    DELETE FROM githashnotes WHERE id = <id_to_delete>;

    DELETE FROM notes WHERE id = <id_to_delete>;

Since the data size is relatively small, the dashboard server can keep a copy of
all recent commit info (e.g., for constructing a "blamelist"), annotations, and
their many-to-many relationship for use in the context.

Password for the database will be stored in the metadata instance. To see the
current password stored in metadata and the fingerprint:

    gcutil --project=google.com:skia-buildbots getinstance [skia-perf GCE instance]

To set the mysql password that perf is to use:

    gcutil --project=google.com:skia-buildbots setinstancemetadata [skia-perf GCE instance] --metadata=readonly:[password-from-valentine] --metadata=readwrite:[password-from-valentine] --fingerprint=[the metadata fingerprint]


Startup and config
------------------
The server is started and stopped via:

    sudo /etc/init.d/perf [start|stop|restart]

But sysv init only handles starting and stopping a program once, so we use
Monit to monitor the application and restart it if it crashes. The config
is in:

    /etc/monit/conf.d/perf

Installation
------------
See the README file.