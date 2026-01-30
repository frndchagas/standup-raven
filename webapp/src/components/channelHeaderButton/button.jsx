import React from 'react';
import PropTypes from 'prop-types';

import logo from '../../assets/images/logo.svg';
import RavenClient from '../../raven-client';

const MAX_APP_BAR_ANCESTOR_DEPTH = 5;

class ChannelHeaderButtonIcon extends React.Component {
    constructor(props) {
        super(props);

        this.myRef = React.createRef();
        this.state = this.getInitialState();
    }

    componentDidMount() {
        RavenClient.Config.getActiveChannels(this.props.siteURL)
            .then((activeChannels) => {
                const activeChannelMap = {};
                activeChannels.forEach((x) => {
                    activeChannelMap[x] = true;
                });
                this.setState({
                    activeChannels: activeChannelMap,
                });
            });
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps.added !== this.props.added || prevProps.removed !== this.props.removed) {
            const activeChannels = this.state.activeChannels;

            if (this.props.added !== prevProps.added) {
                // new active channel is added
                activeChannels[this.props.added] = true;
            }

            if (this.props.removed !== prevProps.removed) {
                // new channel was removed
                activeChannels[this.props.removed] = undefined;
            }

            this.setState({
                activeChannels,
            });
        }
    }

    getInitialState = () => {
        return {
            activeChannels: {},
            parent: undefined,
        };
    };

    handleRef = (ref) => {
        if (ref) {
            this.setState({
                parent: ref.parentNode,
            });
        }
    }

    isChannelHeaderButtonInDropdown = () => {
        try {
            const ancestor = this.state.parent.parentNode.parentNode.parentNode.parentNode;
            if (!ancestor || !ancestor.classList) {
                return false;
            }
            return ancestor.classList.contains('dropdown') && ancestor.classList.contains('btn-group');
        } catch (e) {
            return false;
        }
    }

    getIconParentToHide = () => {
        // In Mattermost v9+ the button may be rendered inside the app bar
        // with a different DOM structure than the old channel header dropdown.
        const appBarParent = this.findAppBarParent();
        if (appBarParent) {
            return appBarParent;
        }

        if (this.isChannelHeaderButtonInDropdown()) {
            return this.state.parent.parentNode.parentNode;
        }
        return this.state.parent;
    }

    findAppBarParent = () => {
        let node = this.state.parent;
        for (let i = 0; i < MAX_APP_BAR_ANCESTOR_DEPTH && node; i++) {
            if (node.id && node.id.startsWith('app-bar-icon-')) {
                return node.querySelector('.app-bar__icon-inner') || node;
            }
            node = node.parentNode;
        }
        return null;
    }

    render() {
        if (this.state.parent) {
            const targetParent = this.getIconParentToHide();
            if (this.state.activeChannels[this.props.channelID]) {
                targetParent.classList.remove('hidden');
            } else {
                targetParent.classList.add('hidden');
            }
        }

        return (
            <span
                ref={this.handleRef}
                style={{
                    width: '1.8em',
                    height: '1.8em',
                }}
                dangerouslySetInnerHTML={{
                    __html: logo,
                }}
            />
        );
    }
}

ChannelHeaderButtonIcon.propTypes = {
    channelID: PropTypes.string.isRequired,
    siteURL: PropTypes.string.isRequired,
    added: PropTypes.string.isRequired,
    removed: PropTypes.string.isRequired,
};

export default ChannelHeaderButtonIcon;
