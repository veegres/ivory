import {Cancel, CheckCircle} from "@mui/icons-material"
import {Box} from "@mui/material"
import {useEffect, useState} from "react"

import {useRouterVaultCreate} from "../../../../api/vault/hook"
import {Vault, VaultType} from "../../../../api/vault/type"
import {SxPropsMap} from "../../../../app/type"
import {VaultOptions} from "../../../../app/utils"
import {SimpleButton} from "../../../view/button/SimpleButton"
import {VaultRow} from "./VaultRow"

const SX: SxPropsMap = {
    box: {paddingRight: "10px", display: "flex", flexDirection: "column", gap: 1},
    head: {marginBottom: "5px"},
    title: {display: "flex", padding: "5px 3px", alignItems: "center", fontSize: 17},
    desc: {fontSize: 12, whiteSpace: "pre-wrap"},
    code: {
        display: "inline-block", border: 1, borderColor: "divider",
        color: "text.secondary", padding: "0px 5px", borderRadius: 1,
    },
}

type Props = {
    type: VaultType,
}

export function VaultItemNew(props: Props) {
    const {type} = props
    const initVault: Vault = {username: "", secret: "", type}
    const options = VaultOptions[type]
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(type, handleCancel)

    useEffect(() => setVault(v => ({...v, type})), [type])

    return (
        <Box sx={SX.box}>
            <Box sx={SX.head}>
                <Box sx={SX.title}>{options.icon}&nbsp;New {options.label}</Box>
                {type === VaultType.SSH_KEY ? (
                    <Box sx={SX.desc}>
                        Follow these steps to add an SSH key to your virtual machine if you haven't set one up yet.
                        <table>
                            <tbody>
                                <tr><td>Generate SSH key</td><td><Box sx={SX.code}>ssh-keygen -t ed25519</Box></td></tr>
                                <tr><td>Authorize the key</td><td><Box sx={SX.code}>{"cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys"}</Box></td></tr>
                                <tr><td>Permissions (optional)</td><td><Box sx={SX.code}>chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys</Box></td></tr>
                            </tbody>
                        </table>
                        Copy your public key that you have created from <Box sx={SX.code}>~/.ssh/id_ed25519.pub</Box> and
                        add it here. All keys are encrypted by the secret key and safely stored inside Ivory. Ivory decrypt them
                        at the moment they're needed.
                    </Box>
                ) : (
                    <Box sx={SX.desc}>
                        All passwords are encrypted by the secret key and safely stored inside Ivory. Ivory decrypt them
                        at the moment they're needed.
                    </Box>
                )}
            </Box>
            <VaultRow
                renderButtons={renderButtons()}
                disabled={false}
                vault={vault}
                onChangeVault={(vault) => {setVault(vault); setClean(true)}}
                onEmpty={(v) => setEmpty(v)}
            />
        </Box>
    )

    function renderButtons() {
        return (
            <>
                <SimpleButton disabled={!clean || createVault.isPending} onClick={handleCancel}><Cancel/></SimpleButton>
                <SimpleButton loading={createVault.isPending} disabled={empty} onClick={handleCreate}><CheckCircle/></SimpleButton>
            </>
        )
    }

    function handleCancel() {
        setVault(initVault)
        setClean(false)
    }

    function handleCreate() {
        createVault.mutate(vault)
    }
}
