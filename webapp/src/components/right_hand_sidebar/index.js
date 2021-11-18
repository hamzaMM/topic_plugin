import {connect} from 'react-redux';

import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {getPostsInCurrentChannel} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    team: getCurrentTeam(state),
    channel: getCurrentChannelId(state),
    posts: getPostsInCurrentChannel(state),
});

export default connect(mapStateToProps)(RHSView);