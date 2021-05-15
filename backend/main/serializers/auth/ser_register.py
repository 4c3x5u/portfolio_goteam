from rest_framework import serializers
import bcrypt
import status

from ...models import User, Team
from ...serializers.user.ser_user import UserSerializer
from main.helpers.custom_api_exception import CustomAPIException
from ...helpers.board_helper import BoardHelper
from ...helpers.tutorial_helper import TutorialHelper


class RegisterSerializer(UserSerializer):
    password_confirmation = serializers.CharField(
        error_messages={
            'blank': 'Password confirmation cannot be empty.'
        }
    )
    invite_code = serializers.UUIDField(
        required=False,
        error_messages={'invalid': 'Invalid invite code.'}
    )

    @staticmethod
    def validate_username(value):
        try:
            User.objects.get(username=value)
        except User.DoesNotExist:
            return value
        raise CustomAPIException('username',
                                 'Username already exists.',
                                 status.HTTP_400_BAD_REQUEST)

    def validate(self, attrs):
        password_confirmation = attrs.pop('password_confirmation')

        if password_confirmation != attrs.get('password'):
            raise CustomAPIException(
                'password_confirmation',
                'Confirmation must match the password.',
                status.HTTP_400_BAD_REQUEST
            )

        invite_code = attrs.get('invite_code')
        if invite_code:
            try:
                attrs['team'] = Team.objects.get(invite_code=invite_code)
            except Team.DoesNotExist:
                raise CustomAPIException('invite_code',
                                         'Team not found.',
                                         status.HTTP_400_BAD_REQUEST)
            attrs['is_admin'] = False
            attrs.pop('invite_code')
        else:
            attrs['is_admin'] = True

        return super().validate(attrs)

    def create(self, validated_data):
        validated_data['password'] = bcrypt.hashpw(
            bytes(validated_data['password'], 'utf-8'),
            bcrypt.gensalt()
        )

        is_admin = validated_data.get('is_admin')
        if is_admin and not validated_data.get('team'):
            validated_data['team'] = Team.objects.create()

        user = User.objects.create(**validated_data)

        if is_admin:
            board_helper = BoardHelper('New Board', user)
            board = board_helper.create_board()
            ready_column = board.column_set.all()[1]

            tutorial_helper = TutorialHelper(user, ready_column)
            tutorial_helper.start()

        return user

    def to_representation(self, instance):
        return {
            'msg': 'Registration successful.',
            'username': instance.username,
            'token': bcrypt.hashpw(
                bytes(instance.username, 'utf-8') + instance.password,
                bcrypt.gensalt()
            ).decode('utf-8'),
            'teamId': instance.team_id,
            'isAdmin': instance.is_admin
        }
