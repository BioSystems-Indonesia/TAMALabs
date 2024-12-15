import {ReferenceInput, ReferenceInputProps} from "react-admin";

export interface FeatureListProps extends Partial<ReferenceInputProps> {
    types: string
    source: string
}

export default function FeatureList(props: FeatureListProps) {
    const {types, source, children, ...rest} = props;

    return (<ReferenceInput {...rest}
                            source={props.source} reference={`feature-list-${types}`}>
        {children}
    </ReferenceInput>)
}