import {AlertButton} from "../../view/button/AlertButton";
import {InstanceRequest} from "../../../type/instance";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";

type Props = {
    request: InstanceRequest,
    cluster: string,
}

export function ReloadButton(props: Props) {
    const {request, cluster} = props

    const options = useMutationOptions([["instance/overview", cluster]])
    const reload = useMutation({mutationFn: instanceApi.reload, ...options})

    return (
        <AlertButton
            color={"info"}
            size={"small"}
            label={"Reload"}
            title={`Make a reload of ${request.sidecar.host}?`}
            description={`It will reload postgres config, it doesn't have any downtime. It won't help if pending 
            restart is true, some parameters require postgres restart.`}
            loading={reload.isPending}
            onClick={() => reload.mutate(request)}
        />
    )
}
