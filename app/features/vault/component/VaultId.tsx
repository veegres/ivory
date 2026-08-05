import {Tooltip} from "@mui/material"

import {getShortUuid} from "../../../shared/helper/HelperUtils"
import {VaultInput} from "./VaultInput"

const SX = {
    input: {width: "90px"},
}

type Props = {
    uuid: string,
}

export function VaultId(props: Props) {
    const {uuid} = props
    return (
        <Tooltip placement={"top-start"} title={uuid}>
            <VaultInput
                sx={SX.input}
                label={"ID"}
                type={"text"}
                value={getShortUuid(uuid)}
                disabled={true}
                onChange={() => {}}
            />
        </Tooltip>
    )
}
