package(default_visibility = ["PUBLIC"])

go_toolchain(
    name = "toolchain",
    hashes = [
        "dab7d9c34361dc21ec237d584590d72500652e7c909bf082758fb63064fca0ef", # linux-amd64
        "b1cb56123f0d9542d82194eae392f94a5712dacfcac46be87d6a9dd79ee8b1df", # darwin-amd64
    ],
    version = "1.17.1",
)

go_module(
    name = "x_sys",
    hashes = ["e9ba045ac49a174a9bb1220fc50e89b46e26c8ea18d987248baec06478c9c7bf"],
    install = [
        "unix",
    ],
    module = "golang.org/x/sys",
    version = "95b1ffbd15a57cc5abb3f04402b9e8ec0016a52c",
)


go_module(
    name = "jessevdk_flags",
    licences = ["BSD-3-Clause"],
    module = "github.com/jessevdk/go-flags",
    version = "v1.5.0",
    deps = [":x_sys"],
)

go_module(
    name = "rs_zerolog",
    install = [
        ".",
        "internal/json",
    ],
    licences = ["MIT"],
    module = "github.com/rs/zerolog",
    version = "v1.23.0",
)
