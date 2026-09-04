import {AlertButton} from "../../../shared/component/button/AlertButton"
import {useStoreAction} from "../../../shared/provider/StoreProvider"
import {useRouterErase} from "../api/ManagementHook"

export function ManageEraseButton() {
    const {clear} = useStoreAction
    const erase = useRouterErase(clear)

    return (
        <AlertButton
            variant={"outlined"}
            color={"error"}
            label={"Erase"}
            title={"Erase all data?"}
            description={"This action will remove all of your data (vaults, certs, logs, queries, etc)."}
            loading={erase.isPending}
            onClick={() => erase.mutate()}
        />
    )
}
