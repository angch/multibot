/*
Copyright © 2021 Ang Chin Han <ang.chin.han@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/angch/discordbot/cmd"
	_ "github.com/angch/discordbot/pkg/apod"
	_ "github.com/angch/discordbot/pkg/askfaz"
	_ "github.com/angch/discordbot/pkg/echo"
	_ "github.com/angch/discordbot/pkg/kulll"
	_ "github.com/angch/discordbot/pkg/meme"
	// _ "github.com/angch/discordbot/pkg/ocr"
	_ "github.com/angch/discordbot/pkg/qrdecode"

	_ "github.com/angch/discordbot/pkg/spacetraders"
	_ "github.com/angch/discordbot/pkg/stablediffusion"

	// _ "github.com/angch/discordbot/pkg/stoic"
	_ "github.com/angch/discordbot/pkg/unicodefont"
	_ "github.com/angch/discordbot/pkg/xkcd"
	_ "github.com/angch/discordbot/pkg/ymca"
	_ "github.com/angch/discordbot/pkg/ynot"
)

func main() {
	cmd.Execute()
}
