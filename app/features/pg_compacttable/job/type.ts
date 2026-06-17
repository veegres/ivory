// COMMON (WEB AND SERVER)

export enum JobStatus {
    PENDING,
    RUNNING,
    FINISHED,
    FAILED,
    STOPPED,
    UNKNOWN,
}

export enum EventType {
    SERVER = "server",
    STATUS = "status",
    LOG = "log",
    STREAM = "stream",
    BROWSER = "browser",
}

export enum EventStreamType {
    START = "start",
    END = "end",
}

// SPECIFIC (WEB)


