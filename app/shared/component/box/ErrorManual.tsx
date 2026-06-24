import {ErrorSmart} from "./ErrorSmart"

export function ErrorMainNodeMissing() {
    return <ErrorSmart error={"Couldn't detect Main node, you have some problems in your setup"}/>
}

export function ErrorLeaderMissing() {
    return <ErrorSmart error={"Node is not a leader, can be used only on leader"}/>
}

export function ErrorDbMissing() {
    return <ErrorSmart error={"Provide Database Port to interact with it"}/>
}

export function ErrorSshMissing() {
    return <ErrorSmart error={"Provide SSH Key and Port to interact with VM"}/>
}

export function ErrorKeeperMissing() {
    return <ErrorSmart error={"Provide Keeper Port to work with it"}/>
}

export function ErrorKeeperRequestMissing() {
    return <ErrorSmart error={"Cannot parse Keeper Request"}/>
}

export function ErrorUserInfoMissing() {
    return <ErrorSmart error={"Can't get user info"}/>
}

export function ErrorNoAccess({name}: {name: string}) {
    return <ErrorSmart error={`No access for ${name} feature, you can request permission in the settings`}/>
}
