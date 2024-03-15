import {AlertButton} from "../../view/button/AlertButton";
import {InstanceRequest} from "../../../type/instance";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {instanceApi} from "../../../app/api";

type Props = {
    request: InstanceRequest,
    cluster: string,
}

export function RestartButton(props: Props) {
    const {request, cluster} = props

    const options = useMutationOptions([["instance/overview", cluster]])
    const restart = useMutation({mutationFn: instanceApi.restart, ...options})

    return (
        <AlertButton
            color={"info"}
            size={"small"}
            label={"Restart"}
            title={`Make a restart of ${request.sidecar.host}?`}
            description={"It will restart postgres, that will cause some downtime."}
            loading={restart.isPending}
            onClick={() => restart.mutate(request)}
        />
    )
}
