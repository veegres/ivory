import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {infoApi} from "../../../app/api";
import {useStore} from "../../../provider/StoreProvider";
import {LoadingButton} from "@mui/lab";
import {useState} from "react";
import {AlertDialog} from "../../view/AlertDialog";

export function EraseButton() {
    const {clear} = useStore()
    const [open, setOpen] = useState(false)

    const cleanOptions = useMutationOptions([["info"]], clear)
    const clean = useMutation(infoApi.erase, cleanOptions)

    return (
        <>
            <LoadingButton
                size={"small"}
                color={"error"}
                variant={"outlined"}
                onClick={() => setOpen(true)}
                loading={clean.isLoading}
            >
                Erase
            </LoadingButton>
            <AlertDialog
                open={open}
                title={"Erase all data?"}
                content={`This action will remove all of your data (passwords, certs, logs, queries, etc).`}
                onAgree={() => clean.mutate()}
                onClose={() => setOpen(false)}
            />
        </>
    )
}
