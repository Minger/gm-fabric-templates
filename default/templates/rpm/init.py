#!/usr/bin/python
#
# The init script for the {{.ServiceName}} service. Save to /etc/init.d
#
#
# chkconfig: - 20 70
# description: {{.ServiceName}}-X.X server process
#
#
### BEGIN INIT INFO
# Provides: {{.ServiceName}}-X.X
# Description: {{.ServiceName}}-X.X server process
### END INIT INFO

import sys, os, subprocess, re, time, signal
from pwd import getpwnam
from grp import getgrnam

USER = "{{.ServiceName}}"
GROUP = "services"

VERSION_MAJOR = "1"
VERSION_MINOR = "1"
VERSION_PATCH = "0"

SERVICE = "{{.ServiceName}}-{0}.{1}".format(VERSION_MAJOR, VERSION_MINOR)
EXE = "/opt/services/{0}/bin/{{.ServiceName}}-bin".format(SERVICE)
CONFIG = "/opt/services/{0}/etc/settings.toml".format(SERVICE)

LOCKFILE = os.path.join('/var/lock/subsys', SERVICE)
PIDFILE = '/var/run/${0}.pid'.format(SERVICE)
LOG_OUT = "/var/log/{0}.log".format(SERVICE)
LOG_ERR = "/var/log/{0}.err".format(SERVICE)


def change_owner(path, uid, gid, perms):
    """
    Set permissions on file or dir. uid and perms must be int type, eg. 0o755.
    Group is ignored with -1.
    """
    statinfo = os.stat(path)
    puid = statinfo.st_uid
    pgid = statinfo.st_gid
    pmode = statinfo.st_mode
    if uid != puid or (gid != -1 and gid != pgid):
        os.chown(path, uid, gid)
    if perms != oct(pmode & 0777):
        os.chmod(path, perms)


def create_dir_not_exists(path):
    if not os.path.exists(path):
        print 'creating %s' % path
        os.makedirs(path)


def lock():
    open(LOCKFILE, 'w').close()


def locked():
    return os.path.exists(LOCKFILE)


def touch_file(path):
    if not os.path.exists(path):
        f = open(path, "w")
        f.close()


def unlock():
    os.remove(LOCKFILE)


def start():
    # Get target uid/gid.
    uid = getpwnam(USER).pw_uid
    gid = getgrnam(GROUP).gr_gid

    # Capture root's environment in dictionary. Pass to Popen.
    env_for_proc = dict()
    for k in os.environ.keys():
        env_for_proc[k] = os.getenv(k, "")

    # Change permissions on configuration file.
    os.chown(CONFIG, uid, -1)

    # Change ownership of logging directory and file, if the file exists.
    logging_dir = os.path.split(LOG_OUT)[0]
    create_dir_not_exists(logging_dir)
    change_owner(logging_dir, uid, gid, 0o750)
    if os.path.exists(LOG_OUT):
        change_owner(LOG_OUT, uid, gid, 0o750)

    # Change ownership of pidfile.
    touch_file(PIDFILE)
    change_owner(PIDFILE, uid, gid, 0o750)

    # Change group/user (note order).
    os.setgid(gid)
    os.setuid(uid)

    print "Starting {0}".format(SERVICE)
    f = open(LOG_OUT, "a")

    command_client = [EXE, '--config', CONFIG]
    ps = subprocess.Popen(command_client, stdout=f, stderr=f, env=env_for_proc)
    with open(PIDFILE, "w") as pf:
        pf.write(str(ps.pid) + '\n')


def stop():
    print "Stopping {0}".format(SERVICE)
    if not os.path.exists(PIDFILE):
        print "error: could not find pidfile at {0}. {1} may be running under a custom script".format(PIDFILE, SERVICE)
        sys.exit(1)
    try:
        # Remove invalid pidfile.
        if not is_service_running():
            print "removing stale pidfile"
            os.remove(PIDFILE)
            sys.exit(1)
        pid = get_pid_from_pidfile(PIDFILE)
        os.kill(int(pid), signal.SIGKILL)
    except Exception as e:
        print e
        print "could not kill pid {0}".format(PIDFILE)
        sys.exit(1)
    os.remove(PIDFILE)


def restart():
    stop()
    lock()
    start()


def status():
    if not locked():
        print "STATUS: {0} not running".format(SERVICE)
        return 3
    else:
        if not os.path.exists(PIDFILE):
            print "STATUS: {0} is not running".format(SERVICE)
            return 0
        if not is_service_running():
            print "pidfile exists but {0} not running".format(SERVICE)
            print "removing stale pidfile"
            os.remove(PIDFILE)
            return 1
        print "STATUS: {0} is running".format(SERVICE)
        return 0


def get_pid_from_pidfile(path):
    with open(path, 'r') as f:
        pid = f.readline().split('\n')[0]
        if len(pid) > 0:
            return pid
        else:
            raise ValueError("invalid pid read from file: %s" % pid)


def is_process_running(process_id):
    try:
        os.kill(int(process_id), 0)
        return True
    except OSError:
        return False


def is_service_running():
    if not os.path.exists(PIDFILE):
        return False
    else:
        pid = get_pid_from_pidfile(PIDFILE)
    return is_process_running(pid)


def is_absolute_path(path):
    return os.path.isabs(path)


def shell_source(script):
    pipe = subprocess.Popen(". %s; env" % script, stdout=subprocess.PIPE, shell=True)
    output = pipe.communicate()[0]
    env = dict((line.split("=", 1) for line in output.splitlines()))
    os.environ.update(env)


# Script entry point.
if __name__ == '__main__':
    try:
        if len(sys.argv) == 1:
            raise ValueError
        create_dir_not_exists(os.path.split(LOG_OUT)[0])
        command = str(sys.argv[1]).strip().lower()
        if command == 'start':
            if not is_service_running():
                lock()
                start()
            sys.exit(0)
        elif command == 'stop':
            stop()
            unlock()
            sys.exit(0)
        elif command == 'restart' or command == 'force-reload':
            restart()
            sys.exit(0)
        elif command == 'status':
            ok = status()
            sys.exit(ok)
        else:
            raise ValueError
    except (SystemExit):
        pass
    except (ValueError):
        print >> sys.stderr, "Usage: {0} [start|stop|restart|status]".format(SERVICE)
sys.exit(2)
