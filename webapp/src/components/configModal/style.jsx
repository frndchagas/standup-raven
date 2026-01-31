function getStyle() {
    return {
        controlLabel: {
            paddingRight: '10px',
            width: '180px',
            flexShrink: 0,
        },
        controlLabelX: {
            paddingRight: '10px',
            paddingLeft: '10px',
            display: 'inline-flex',
            alignItems: 'center',
            minHeight: '40px',
        },
        formField: {
            flex: 1,
            minWidth: 0,
        },
        windowTimeRow: {
            display: 'inline-flex',
            alignItems: 'center',
            flexWrap: 'wrap',
            gap: '4px',
        },
        formGroup: {
            marginBottom: '16px',
            minHeight: '40px',
            display: 'flex',
            flexWrap: 'wrap',
            alignItems: 'center',
        },
        formGroupNoMarginBottom: {
            marginBottom: '0',
        },
        sections: {
            marginBottom: '10px',
        },
        sectionGroup: {
            maxHeight: '300px',
            overflowY: 'auto',
        },
        spinner: {
            width: '80px',
            display: 'block',
            margin: '50px auto',
        },
        scrollY: {
            overflowY: 'scroll',
        },
        alert: {
            width: '90%',
            marginLeft: 'auto',
            marginRight: 'auto',
            textAlign: 'center',
            borderRadius: '5px',
            whiteSpace: 'pre-line',
            animation: 'pop 0.2s ease-in',
        },
        body: {
            minHeight: '380px',
        },
        bodyCompact: {
            minHeight: 'unset',
        },
        standupErrorSection: {
            textAlign: 'center',
            color: 'var(--center-channel-color)',
        },
    };
}

export default {
    getStyle,
};
