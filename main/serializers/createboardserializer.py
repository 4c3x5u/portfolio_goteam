from rest_framework import serializers
from ..models import User, Team, Board


class CreateBoardSerializer(serializers.Serializer):
    username = serializers.CharField(min_length=5, max_length=35)
    board_team = None

    def update(self, instance, validated_data):
        raise Exception('Update not allowed on this endpoint.')

    def validate(self, data):
        try:
            user = User.objects.get(username=data.get('username'))
        except User.DoesNotExist:
            raise serializers.ValidationError({
                'username': "Invalid username."
            })
        if not user.is_admin:
            raise serializers.ValidationError(
                'Only the team admin can create a board.',
                'not_authorized'
            )
        self.board_team = Team.objects.get(id=user.team.id)
        return super().validate(data)

    def create(self, validated_data):
        return Board.objects.create(team=self.board_team)

