import {AlertButton} from "../../../shared/component/button/AlertButton"
import {useStoreAction} from "../../../shared/provider/StoreProvider"
import {useRouterEraseInitial, useRouterEraseSafe} from "../api/ManagementHook"

type Props = {
    safe: boolean,
}

export function ManageEraseButton(props: Props) {
    const {clear} = useStoreAction
    const cleanInitial = useRouterEraseInitial(clear)
    const cleanSafe = useRouterEraseSafe(clear)

    return (
        <AlertButton
            variant={"outlined"}
            size={"small"}
            color={"error"}
            label={"Erase"}
            title={"Erase all data?"}
            description={"This action will remove all of your data (vaults, certs, logs, queries, etc)."}
            loading={cleanSafe.isPending || cleanInitial.isPending}
            onClick={() => props.safe ? cleanSafe.mutate() : cleanInitial.mutate()}
        />
    )
}
