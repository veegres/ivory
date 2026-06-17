import {AddCircle, Cancel} from "@mui/icons-material"
import {useState} from "react"

import {Code} from "../../../shared/component/box/Code"
import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {useRouterVaultCreate} from "../api/hook"
import {VaultType} from "../api/type"
import {VaultNewWrapper} from "./VaultNewWrapper"
import {VaultRowSshKey} from "./VaultRowSshKey"

export function VaultNewSshKey() {
    const type = VaultType.SSH_KEY
    const initVault = {username: "", secret: "", type}
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(type, handleCancel)

    return (
        <VaultNewWrapper description={renderDescription()}>
            <VaultRowSshKey
                renderButtons={renderButtons()}
                disabled={false}
                vault={vault}
                onChangeVault={(vault) => {setVault(vault); setClean(true)}}
                onEmpty={(v) => setEmpty(v)}
            />
        </VaultNewWrapper>
    )

    function renderDescription() {
        return (
            <>
                Ivory will generate a secure SSH key pair for you.
                The private key will be safely stored in the vault, and you will be able to copy the public key
                to add it to your virtual machine's <Code>~/.ssh/authorized_keys</Code> file.
            </>
        )
    }

    function renderButtons() {
        return (
            <>
                <SimpleButton disabled={!clean || createVault.isPending} onClick={handleCancel}><Cancel/></SimpleButton>
                <SimpleButton loading={createVault.isPending} disabled={empty} onClick={handleCreate}>
                    <AddCircle/>
                </SimpleButton>
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
