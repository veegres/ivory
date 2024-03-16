import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {useStoreAction} from "../../../provider/StoreProvider";
import {InitialApi, SafeApi} from "../../../app/api";
import {AlertButton} from "../../view/button/AlertButton";

type Props = {
    safe: boolean,
}

export function EraseButton(props: Props) {
    const {clear} = useStoreAction()

    const cleanOptions = useMutationOptions([["info"]], clear)
    const cleanInitial = useMutation({mutationFn: InitialApi.erase, ...cleanOptions})
    const cleanSafe = useMutation({mutationFn: SafeApi.erase, ...cleanOptions})

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
