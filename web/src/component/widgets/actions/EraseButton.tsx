import {useStoreAction} from "../../../provider/StoreProvider";
import {AlertButton} from "../../view/button/AlertButton";
import {useRouterEraseInitial, useRouterEraseSafe} from "../../../api/management/hook";

type Props = {
    safe: boolean,
}

export function EraseButton(props: Props) {
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
            description={`This action will remove all of your data (passwords, certs, logs, queries, etc).`}
            loading={cleanSafe.isPending || cleanInitial.isPending}
            onClick={() => props.safe ? cleanSafe.mutate() : cleanInitial.mutate()}
        />
    )
}
