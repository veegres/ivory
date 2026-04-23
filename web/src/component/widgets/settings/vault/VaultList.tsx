import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {Vault, VaultMap} from "../../../../api/vault/type"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {NoBox} from "../../../view/box/NoBox"
import {VaultItem} from "./VaultItem"

type Props = {
    vaults?: VaultMap,
    error: any,
}

export function VaultList(props: Props) {
    const {vaults, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Vault>(vaults ?? {})
    if (list.length === 0) return <NoBox text={"There are no vaults yet"}/>

    return (
        <TransitionGroup appear={false}>
            {list.map(([key, vault]) => (
                <Collapse key={key}>
                    <VaultItem uuid={key} vault={vault}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
