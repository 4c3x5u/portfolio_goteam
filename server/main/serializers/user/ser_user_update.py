from rest_framework import serializers
import status

from server.main.serializers.user.ser_user import UserSerializer
from server.main.helpers.auth_helper import AuthHelper
from server.main.helpers.custom_api_exception import CustomAPIException
from server.main.models import User, Board


class UpdateUserSerializer(UserSerializer):
    """
    Used only for adding/removing a user to/from a board
    """
    board = serializers.PrimaryKeyRelatedField(
        queryset=Board.objects.prefetch_related('user', 'team__user_set')
                              .all(),
        error_messages={'blank': 'Board ID cannot be blank.',
                        'null': 'Board ID cannot be null.',
                        'incorrect_type': 'Board ID must be a number.',
                        'does_not_exist': 'Board does not exist.'}
    )
    is_active = serializers.BooleanField(error_messages={
        'blank': 'Is Active cannot be blank.',
        'null': 'Is Active cannot be null.',
        'invalid': 'Is Active must be a boolean.'
    })
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = UserSerializer.Meta.model
        fields = 'username', 'board', 'is_active', 'auth_user', \
                 'auth_token',

    def validate(self, attrs):
        authenticated_user = AuthHelper.authenticate(attrs.pop('auth_user'),
                                                     attrs.pop('auth_token'))

        board = attrs.get('board')

        try:
            user = board.team.user_set.get(username=attrs.get('username'))
        except User.DoesNotExist:
            raise CustomAPIException('username',
                                     'User not found.',
                                     status.HTTP_404_NOT_FOUND)

        AuthHelper.authorize(authenticated_user, user.team_id)

        self.instance = {'user': user,
                         'board': board,
                         'is_active': attrs.get('is_active')}
        return attrs

    def update(self, instance, validated_data):
        board = instance.get('board')
        user = instance.get('user')
        if instance.get('is_active'):
            board.user.add(user)
        else:
            board.user.remove(user)
        return instance

    def to_representation(self, instance):
        user = instance.get('user')
        board = instance.get('board')
        action = 'added' if instance.get('is_active') else 'removed'
        return {
            'msg': f'{user.username} is {action} from {board.name}.'
        }

