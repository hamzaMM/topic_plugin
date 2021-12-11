import {connect} from 'react-redux';

import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';
import {getPostsInCurrentChannel} from 'mattermost-redux/selectors/entities/posts';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    user: getCurrentUserId(state),
    channel: getChannel(state, '7hwrm5ajdiy69mf3oue1cich6r'),
    posts: getPostsInCurrentChannel(state),
});

export default connect(mapStateToProps)(RHSView);