function getStyle() {
    return {
        header: {
            marginTop: '0',
            marginBottom: '30px',
        },
        controlBtns: {
            marginRight: '8px',
            minWidth: '40px',
            minHeight: '40px',
        },
        controlBar: {
            display: 'flex',
            gap: '8px',
            flexWrap: 'wrap',
            alignItems: 'center',
        },
        form: {
            minHeight: '200px',
            maxHeight: '50vh',
            overflowY: 'auto',
            marginBottom: '16px',
        },
        alert: {
            width: '90%',
            marginLeft: 'auto',
            marginRight: 'auto',
            textAlign: 'center',
            borderRadius: 'var(--radius-s, 4px)',
            whiteSpace: 'pre-line',
            animation: 'pop 0.2s ease-in',
        },
        spinner: {
            width: '80px',
            display: 'block',
            margin: '50px auto',
        },
        standupErrorMessage: {
            fontWeight: 'bold',
        },
    };
}

export default {
    getStyle,
};
