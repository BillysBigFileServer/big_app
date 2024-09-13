package preferences

import (
	"context"

	"fyne.io/fyne/v2"
)

type preferencesContextKeyType struct{}

var preferencesContextKey preferencesContextKeyType

func ContextWithPreferences(ctx context.Context, preferences fyne.Preferences) context.Context {
	return context.WithValue(ctx, preferencesContextKey, preferences)
}

func GetPreferences(ctx context.Context) fyne.Preferences {
	return ctx.Value(preferencesContextKey).(fyne.Preferences)
}
