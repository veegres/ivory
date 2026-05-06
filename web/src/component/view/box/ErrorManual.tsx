import {ErrorSmart} from "./ErrorSmart"

export function ErrorMainNodeMissing() {
    return <ErrorSmart error={"Couldn't detect Main node, you have some problems in your setup"}/>
}

export function ErrorLeaderMissing() {
    return <ErrorSmart error={"Main node is not a leader, probably something has happened or you've change it"}/>
}

export function ErrorDbMissing() {
    return <ErrorSmart error={"Provide Database Port to interact with it"}/>
}

export function ErrorSshMissing() {
    return <ErrorSmart error={"Provide SSH Key and Port to see interact with VM"}/>
}

export function ErrorKeeperMissing() {
    return <ErrorSmart error={"Provide Keeper Port to work with it"}/>
}

export function ErrorUserInfoMissing() {
    return <ErrorSmart error={"Can't get user info"}/>
}
