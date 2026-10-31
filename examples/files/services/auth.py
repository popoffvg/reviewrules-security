import logging

logger = logging.getLogger(__name__)


def login(store, username, password):
    user = store.check(username, password)
    if user is None:
        logger.warning("login failed user=%s password=%s", username, password)
        return None
    logger.info("login ok user_id=%s", user.id)
    return user
