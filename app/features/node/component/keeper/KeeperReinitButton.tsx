import {Box, FormControlLabel, Switch} from "@mui/material"
import {useState} from "react"

import {AlertButton} from "../../../../shared/component/button/AlertButton"
import {SxPropsMap} from "../../../../shared/helper/type"
import {Feature} from "../../../feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeReinit} from "../../api/hook"
import {KeeperOneRequest} from "../../api/type"

const SX: SxPropsMap = {
    force: {margin: "0px"},
    description: {color: "text.disabled", fontSize: "12px"},
}

type Props = {
    request: KeeperOneRequest,
    cluster: string,
}

export function KeeperReinitButton(props: Props) {
    const {request, cluster} = props
    const [force, setForce] = useState(false)
    const reinit = useRouterNodeReinit(cluster)

    const body = {force}

    return (
        <ManageAccess feature={Feature.ManageNodeDbReinitialize}>
            <AlertButton
                color={"info"}
                size={"small"}
                label={"Reinit"}
                title={`Make a reinit of ${request.host}?`}
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
        </ManageAccess>
    )

    function renderLabel() {
        return (
            <Box>
                <Box>Force</Box>
                <Box sx={SX.description}>
                    In order to overcome fail if Patroni is in a loop trying to recover (restart) a failed Postgres.
                </Box>
            </Box>
        )
    }

    function handleClick() {
        reinit.mutate({...request, body})
        setForce(false)
    }
}
