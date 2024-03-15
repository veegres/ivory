import {AlertButton} from "../../view/button/AlertButton";
import {InstanceRequest} from "../../../type/instance";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";

type Props = {
    request: InstanceRequest,
    cluster: string,
}

export function ReinitButton(props: Props) {
    const {request, cluster} = props

    const options = useMutationOptions([["instance/overview", cluster]])
    const reinit = useMutation({mutationFn: instanceApi.reinitialize, ...options})

    return (
        <AlertButton
            color={"primary"}
            label={"Reinit"}
            title={`Make a reinit of ${request.sidecar.host}?`}
            description={"It will erase all node data and will download it from scratch."}
            loading={reinit.isPending}
            onClick={() => reinit.mutate(request)}
        />
    )
}
