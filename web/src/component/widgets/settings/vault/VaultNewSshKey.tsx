import {AddCircle, Cancel} from "@mui/icons-material"
import {useState} from "react"

import {useRouterVaultCreate} from "../../../../api/vault/hook"
import {VaultType} from "../../../../api/vault/type"
import {VaultOptions} from "../../../../app/utils"
import {Code} from "../../../view/box/Code"
import {SimpleButton} from "../../../view/button/SimpleButton"
import {VaultNewWrapper} from "./VaultNewWrapper"
import {VaultRowSshKey} from "./VaultRowSshKey"

export function VaultNewSshKey() {
    const type = VaultType.SSH_KEY
    const initVault = {username: "", secret: "", type}
    const options = VaultOptions[type]
    const [vault, setVault] = useState(initVault)
    const [empty, setEmpty] = useState(false)
    const [clean, setClean] = useState(false)
    const createVault = useRouterVaultCreate(type, handleCancel)

    return (
        <VaultNewWrapper
            icon={options.icon}
            label={options.label}
            description={renderDescription()}
        >
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
