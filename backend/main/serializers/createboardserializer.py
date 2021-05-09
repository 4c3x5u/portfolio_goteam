from rest_framework import serializers
import status

from .boardserializer import BoardSerializer
from ..validation.val_auth import authenticate_custom, authorization_error, \
    authorize_custom
from ..validation.val_custom import CustomAPIException
from ..util import create_board


class CreateBoardSerializer(BoardSerializer):
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta(BoardSerializer.Meta):
        fields = 'auth_user', 'auth_token', 'team', 'name'

    @staticmethod
    def validate_team_id(value):
        if not value:
            raise CustomAPIException('team_id',
                                     'Team ID cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)
        try:
            int(value)
        except ValueError:
            raise CustomAPIException('team_id',
                                     'Team ID must be a number.',
                                     status.HTTP_400_BAD_REQUEST)
        return value

    def validate(self, attrs):
        auth_user = attrs.pop('auth_user')
        auth_token = attrs.pop('auth_token')
        user, authentication_error = authenticate_custom(auth_user, auth_token)
        if authentication_error:
            raise authentication_error

        authorize_error = authorize_custom(user.username)
        if authorize_error:
            raise authorize_error

        team = attrs.get('team')
        if team.id != user.team_id:
            raise authorization_error

        board_name = attrs.get('name')
        board_serializer = BoardSerializer(data={
            'team': team.id,
            'name': board_name,
        })
        if not board_serializer.is_valid():
            raise board_serializer.errors

        return {'board_name': board_name,
                'team_id': team.id,
                'team_admin': team.user_set.get(username=user.username)}

    def create(self, validated_data):
        board, create_error = create_board(
            name=validated_data.get('board_name'),
            team_id=validated_data.get('team_id'),
            team_admin=validated_data.get('team_admin')
        )
        if create_error:
            raise create_error
        return board

    def to_representation(self, instance):
        return {'msg': 'Board creation successful.',
                'id': instance.id}
