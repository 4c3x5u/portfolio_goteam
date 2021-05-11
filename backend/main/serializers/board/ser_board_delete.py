from rest_framework import serializers

from main.serializers.board.ser_board import BoardSerializer
from main.validation.val_auth import authenticate, authorize
from main.models import Board


class DeleteBoardSerializer(serializers.ModelSerializer):
    board = serializers.PrimaryKeyRelatedField(
        queryset=Board.objects.all(),
        error_messages={'null': 'Board ID cannot be null.',
                        'incorrect_type': 'Board ID must be a number.',
                        'does_not_exist': 'Board does not exist.'}
    )
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = BoardSerializer.Meta.model
        fields = 'board', 'auth_user', 'auth_token',

    def validate(self, attrs):
        user = authenticate(attrs.get('auth_user'), attrs.get('auth_token'))
        board = attrs.get('board')
        authorize(user, board.team_id)
        return board

    def delete(self):
        self.instance = {'id': self.validated_data.id}
        return self.validated_data.delete()

    def to_representation(self, instance):
        return {
            'msg': 'Board deleted successfully.',
            'id': instance.get('id'),
        }



