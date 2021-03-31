from rest_framework import serializers
from ..models import User, Team, Board


class BoardSerializer(serializers.Serializer):
    username = serializers.CharField(min_length=5, max_length=35)
    board_team = None

    def validate(self, data):
        try:
            user = User.objects.get(username=data.get('username'))
        except User.DoesNotExist:
            raise serializers.ValidationError({
                'username': "Invalid username."
            })
        if not user.is_admin:
            raise serializers.ValidationError({
                'user_not_admin': 'Only the team admin can create a board.'
            })
        try:
            team = Team.objects.get(id=user.team.id)
        except Team.DoesNotExist:
            raise serializers.ValidationError({
                'team_id': 'Invalid team ID.'
            })
        self.board_team = team
        return super().validate(data)

    def create(self, validated_data):
        return Board.objects.create(team=self.board_team)

    def update(self, instance, validated_data):
        pass

