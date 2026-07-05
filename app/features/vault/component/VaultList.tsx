import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {NoBox} from "../../../shared/component/box/NoBox"
import {Vault, VaultMap} from "../api/VaultType"
import {VaultListItem} from "./VaultListItem"

type Props = {
    vaults?: VaultMap,
    error: any,
}

export function VaultList(props: Props) {
    const {vaults, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Vault>(vaults ?? {})
    if (list.length === 0) return <NoBox text={"There are no Credentials yet"}/>

    return (
        <TransitionGroup appear={false}>
            {list.map(([key, vault]) => (
                <Collapse key={key}>
                    <VaultListItem uuid={key} vault={vault}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
