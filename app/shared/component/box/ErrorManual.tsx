import {ErrorSmart} from "./ErrorSmart"

export function ErrorMainNodeMissing() {
    return <ErrorSmart error={"Couldn't detect Main node, you have some problems in your setup"}/>
}

export function ErrorLeaderMissing() {
    return <ErrorSmart error={"Node is not a leader, can be used only on leader"}/>
}

export function ErrorDbMissing() {
    return <ErrorSmart error={"Provide Database Port to send queries to the database or interact with database tools"}/>
}

export function ErrorSshMissing() {
    return <ErrorSmart error={"Provide SSH Key and VM Port to interact with your system"}/>
}

export function ErrorKeeperMissing() {
    return <ErrorSmart error={"Provide Keeper Port to be able to use features like switchover, reload, reinit, config, etc"}/>
}

export function ErrorUserInfoMissing() {
    return <ErrorSmart error={"Can't get user info"}/>
}

export function ErrorNoAccess({name}: {name: string}) {
    return <ErrorSmart error={`No access for ${name} feature, you can request permission in the settings`}/>
}

export function ErrorNotSupported({name}: {name: string}) {
    return <ErrorSmart type={"info"} error={`The ${name} feature is not supported by the current cluster's plugin`}/>
}
