import {Options} from "../../../core/widgets/options/Options"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {Options as ClusterOptions} from "../api/ClusterType"

type Props = {
    options: ClusterOptions,
    onUpdate: (options: ClusterOptions) => void,
}

// ClusterOptionsBox is the cluster options widget as a dialog shows it: folded
// away under its own heading, with the plugin selectors disabled, since inside
// a dialog the plugins come from the cluster list filter rather than the form.
export function ClusterOptionsBox(props: Props) {
    const {options, onUpdate} = props

    return (
        <TitleBox label={"All Options"} dense={true}>
            <Options options={options} onUpdate={onUpdate} disablePlugins={true}/>
        </TitleBox>
    )
}
