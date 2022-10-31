from rest_framework import serializers
from rest_framework.exceptions import ValidationError

from server.main.serializers.board.ser_board import BoardSerializer
from server.main.helpers.auth_helper import AuthHelper
from server.main.helpers.board_helper import BoardHelper


class CreateBoardSerializer(BoardSerializer):
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta(BoardSerializer.Meta):
        fields = 'name', 'team', 'auth_user', 'auth_token'

    def validate(self, attrs):
        user = AuthHelper.authenticate(attrs.get('auth_user'),
                                       attrs.get('auth_token'))

        team = attrs.get('team')
        AuthHelper.authorize(user, team.id)

        board_name = attrs.get('name')
        board_serializer = BoardSerializer(data={
            'team': team.id,
            'name': board_name,
        })
        if not board_serializer.is_valid():
            raise ValidationError({'board': board_serializer.errors})

        return {'board_name': board_name,
                'team_id': team.id,
                'team_admin': team.user_set.get(username=user.username)}

    def create(self, validated_data):
        board_helper = BoardHelper(validated_data.get('board_name'),
                                   validated_data.get('team_admin'))
        return board_helper.create_board()

    def to_representation(self, instance):
        return {'msg': 'Board creation successful.',
                'id': instance.id}
