import {AlertButton} from "../../view/button/AlertButton";
import {InstanceRequest} from "../../../type/instance";
import {useMutationOptions} from "../../../hook/QueryCustom";
import {useMutation} from "@tanstack/react-query";
import {InstanceApi} from "../../../app/api";
import {useState} from "react";
import {Box, FormControlLabel, Switch} from "@mui/material";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    force: {margin: "0px"},
    description: {color: "text.disabled", fontSize: "12px"},
}

type Props = {
    request: InstanceRequest,
    cluster: string,
}

export function ReinitButton(props: Props) {
    const {request, cluster} = props

    const [force, setForce] = useState(false)

    const options = useMutationOptions([["instance/overview", cluster]])
    const reinit = useMutation({mutationFn: InstanceApi.reinitialize, ...options})

    const body = {force}

    return (
        <AlertButton
            color={"info"}
            label={"Reinit"}
            title={`Make a reinit of ${request.sidecar.host}?`}
            description={"It will erase all node data and will download it from scratch."}
            loading={reinit.isPending}
            onClick={handleClick}
        >
            <FormControlLabel
                sx={SX.force}
                labelPlacement={"start"}
                control={<Switch checked={force} onClick={() => setForce(!force)}/>}
                label={renderLabel()}
            />
        </AlertButton>
    )

    function renderLabel() {
        return (
            <Box>
                <Box>Force</Box>
                <Box sx={SX.description}>
                    in order to overcome fail if Patroni is in a loop trying to recover (restart) a failed Postgres
                </Box>
            </Box>
        )
    }

    function handleClick() {
        reinit.mutate({...request, body})
        setForce(false)
    }
}
