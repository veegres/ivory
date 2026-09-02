import {DialogScreen} from "./DialogScreen"
import {Logs} from "./Logs"

// NOTE: the virtualizer inside Logs needs a pixel height, so this is where it
// and the dialog have to agree: --size-dialog (600) less the screen's padding
// (10 top, 10 bottom), less the box's own frame (its padding, heading and
// gap), less the logs' own toolbar and gap. Under rather than exactly on it - a
// few px of slack costs nothing, overshooting brings back the outer scrollbar
// this screen exists without.
const LOGS_HEIGHT = 493

type Props = {
    logs: string[],
}

// DialogLogsScreen is a whole dialog screen filled by the output of whatever
// the screen before it ran. It knows nothing about what produced the lines -
// the height is the only thing it decides, and that is the dialog's geometry
// rather than anything about the caller.
export function DialogLogsScreen(props: Props) {
    const {logs} = props

    return (
        <DialogScreen fit={true}>
            <Logs logs={logs} height={LOGS_HEIGHT} auto={false}/>
        </DialogScreen>
    )
}
