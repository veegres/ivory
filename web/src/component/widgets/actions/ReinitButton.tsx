import {Box, FormControlLabel, Switch} from "@mui/material"
import {useState} from "react"

import {Feature} from "../../../api/feature"
import {useRouterNodeReinit} from "../../../api/node/hook"
import {KeeperRequest} from "../../../api/node/type"
import {SxPropsMap} from "../../../app/type"
import {AlertButton} from "../../view/button/AlertButton"
import {Access} from "../access/Access"

const SX: SxPropsMap = {
    force: {margin: "0px"},
    description: {color: "text.disabled", fontSize: "12px"},
}

type Props = {
    request: KeeperRequest,
    cluster: string,
}

export function ReinitButton(props: Props) {
    const {request, cluster} = props
    const [force, setForce] = useState(false)
    const reinit = useRouterNodeReinit(cluster)

    const body = {force}

    return (
        <Access feature={Feature.ManageNodeDbReinitialize}>
            <AlertButton
                color={"info"}
                label={"Reinit"}
                title={`Make a reinit of ${request.connection.host}?`}
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
        </Access>
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
